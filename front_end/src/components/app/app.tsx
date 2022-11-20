import "bootstrap/dist/css/bootstrap.min.css";
import React, { useEffect, useState } from "react";
import "./App.css";
import { CAMERAS } from "../../hasura/camera";
import { useSubscription } from "@apollo/client";
import { Container, Row } from "react-bootstrap";
import { EVENT_DATES, getDeduplicatedDates } from "../../hasura/event_date";
import { Menu } from "./Menu";
import { TYPES } from "../../hasura/type";
import { Content } from "./Content";

function App() {
  const camerasQuery = useSubscription(CAMERAS);
  const datesQuery = useSubscription(EVENT_DATES);
  const types = TYPES;

  const deduplicatedDates = getDeduplicatedDates(datesQuery?.data);

  const [camera, setCamera] = useState(null);
  const [date, setDate] = useState(null);
  const [type, setType] = useState(null);
  const [responsive, setResponsive] = useState(window.innerWidth < 992);

  useEffect(() => {
    const handleResize = () => {
      setResponsive(window.innerWidth < 992);
    };

    window.addEventListener("resize", () => {
      handleResize();
    });

    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  if (!camera && camerasQuery?.data?.camera) {
    setCamera(camerasQuery?.data?.camera[0]);
  }

  if (!date && deduplicatedDates?.length) {
    setDate(deduplicatedDates[0] as any);
  }

  if (!type && types?.length) {
    setType(types[0] as any);
  }

  return (
    <Container style={{ width: "100%", height: "100%" }}>
      <Row>
        <Menu
          cameras={camerasQuery?.data?.camera}
          camera={camera}
          setCamera={(x) => {
            setCamera(x);
          }}
          dates={deduplicatedDates}
          date={date}
          setDate={(x) => {
            setDate(x);
          }}
          types={types}
          type={type}
          setType={(x) => {
            setType(x);
          }}
        />
      </Row>
      <Row>
        <Content
          camera={camera}
          date={date}
          type={type}
          responsive={responsive}
        />
      </Row>
    </Container>
  );
}

export default App;
