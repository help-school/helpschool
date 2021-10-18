import React from 'react'
import {useAuth0} from "@auth0/auth0-react";
import {Col, Container, Row } from "react-bootstrap";
import useProtectedApi from "../hooks/api"

function DonationsList({className = ''} = {}) {
  const {isAuthenticated} = useAuth0();
  const { loading, error, donations } = useProtectedApi('/api/my-donations');

  if (!isAuthenticated) return <div>Please log in to see your donations</div>

  if (loading) return <div>Loading donations...</div>
  if (error) return <div>Failed to load donations: {error.message}</div>
  if (donations.length === 0) return <div>No donations yet</div>

  return <code>{JSON.stringify(donations)}</code>
}

function MyDonations() {
  return (
      <section className={"ftco-section"}>
          <Container>
              <Row>
                  <Col>
                    <h2>Your donations</h2>
                    <DonationsList/>
                  </Col>
              </Row>
          </Container>
      </section>)
}

export default MyDonations