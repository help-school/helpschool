import React from 'react';
import {Nav} from "react-bootstrap";
import {useAuth0} from "@auth0/auth0-react";

function Login() {
    const {
        loginWithRedirect, 
        isAuthenticated, 
        isLoading, 
        logout: logoutWithRedirect
    } = useAuth0();

    function login(key, e) {
        e.preventDefault()
        loginWithRedirect()
    }
    
    function logout(key, e) {
        e.preventDefault()
        logoutWithRedirect({})
    }

    if (isLoading) return <Nav.Link>Logging in...</Nav.Link>
    
    if (isAuthenticated) {
        return <Nav.Link href="/logout" onSelect={logout}>Log out</Nav.Link>
    }
    
    return <Nav.Link href="/login" onSelect={login}>Login</Nav.Link>
}

export default Login;
