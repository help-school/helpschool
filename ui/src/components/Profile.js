import {useAuth0} from "@auth0/auth0-react";
import useApi from "../hooks/api"

function Profile({className = ''} = {}) {
  const {user, isAuthenticated} = useAuth0();
  const { loading, error, data } = useApi('/secret');
  
  if (!isAuthenticated) return <div>Not logged in</div>
  if (loading) return <div>Loading scret...</div>
  if (error) return <div>Error: {error.message}</div>

  return <div className={className}>Hi {user.name}, the secrte is {data.secret}</div>
}

export default Profile