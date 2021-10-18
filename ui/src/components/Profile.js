import React from 'react'
import {useAuth0} from "@auth0/auth0-react";

function Profile({className = ''} = {}) {
  const {user, isAuthenticated} = useAuth0();
  if (!isAuthenticated) return null
  return <div className={className}>Hi {user.name}</div>
}

export default Profile