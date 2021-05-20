import React, {useRef,  useState} from 'react';
import {Col, Container, Form, Button} from "react-bootstrap";
import Alert from 'react-bootstrap/Alert'
import {useDispatch,useSelector} from 'react-redux';
import * as Actions from '../store/actions';

import Introduce from "./Introduce";

function Login() {

  const dispatch = useDispatch();

  const [status, setStatus] = useState("");

  
  
  function handleSubmit(event){
     
      event.preventDefault();
      //console.log(item)  
      // dispatch(Actions.postTeacherRequest(item));
      // setStatus("success")
      // setStatus("failure")

      console.log("Submit clicked: ", event.target)
  }
   
  return (
          <section className={"ftco-section"}>
            <Container>
              <Alert.Heading>Not logged in yet ! </Alert.Heading>
            </Container>
        </section>
    );
   
  }

export default Login;
