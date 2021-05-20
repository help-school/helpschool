import React, {useRef,  useState} from 'react';
import {Col, Container, Form, Button} from "react-bootstrap";
import Alert from 'react-bootstrap/Alert'
import {useDispatch,useSelector} from 'react-redux';
import * as Actions from '../store/actions';

import Introduce from "./Introduce";

function Login() {

  const dispatch = useDispatch();

  const [status, setStatus] = useState("");

  const teacherNameRef = useRef(null)
  const teacherPhoneRef = useRef(null)
  const teacherEmailRef = useRef(null)
  const productUrlRef = useRef(null)
  const quantityRef = useRef(null)
  const schoolNameRef = useRef(null)
  const addressRef = useRef(null)
  const districtRef = useRef(null)
  const stateRef = useRef(null)
  const pincodeRef = useRef(null)
  const postTeachersRequest = useSelector(({schoolReducer}) => schoolReducer.postTeachersRequest);

  // function handleChange(event) {
  //   alert(event.target.name + ":" + event.target.value )
  //   alert(teacherNameRef.current.value)
  //   setTeacherName( event.target.value)
  // }

  function handleSubmit(event){
     
      event.preventDefault();
      const item = {
        teacher_name: teacherNameRef.current.value,
        teacher_email: teacherEmailRef.current.value,
        teacher_phone: teacherPhoneRef.current.value,
        url: productUrlRef.current.value,
        quantity_needed: quantityRef.current.value,
        address: schoolNameRef.current.value,
        place: addressRef.current.value,
        district: districtRef.current.value,
        state: stateRef.current.value,
        country: "India",
        photo_link: "",
        extra_info: "{}",
        pincode: pincodeRef.current.value}
     
      console.log(item)  
      dispatch(Actions.postTeacherRequest(item));
      setStatus("success")
      setStatus("failure")

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
