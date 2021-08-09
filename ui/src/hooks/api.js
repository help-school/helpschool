import { useEffect, useState } from 'react';
import { useAuth0 } from '@auth0/auth0-react';

// see README/.env
const baseURL = process.env.REACT_APP_API_URL
const defaultAudience = process.env.REACT_APP_API_AUDIENCE
const defaultScope = 'read:current_user'

console.log('baseURL', baseURL, 'audience', defaultAudience)

// see https://github.com/auth0/auth0-react/blob/master/EXAMPLES.md

const useApi = (path, options = {}) => {
  const { getAccessTokenSilently } = useAuth0();
  const [state, setState] = useState({
    error: null,
    loading: true,
    data: null,
  });
  const [refreshIndex, setRefreshIndex] = useState(0);

  path = path.startsWith('/') ? path.slice(1) : path

  useEffect(() => {
    (async () => {
      try {
        const { audience = defaultAudience, scope = defaultScope, ...fetchOptions } = options;
        const accessToken = await getAccessTokenSilently({ audience });
        const res = await fetch(baseURL + '/' + path, {
          ...fetchOptions,
          headers: {
            ...fetchOptions.headers,
            // Add the Authorization header to the existing headers
            Authorization: `Bearer ${accessToken}`,
          },
        });
        setState({
          ...state,
          data: await res.json(),
          error: null,
          loading: false,
        });
      } catch (error) {
        setState({
          ...state,
          error,
          loading: false,
        });
      }
    })();
  }, [refreshIndex]);

  return {
    ...state,
    refresh: () => setRefreshIndex(refreshIndex + 1),
  };
};

export default useApi