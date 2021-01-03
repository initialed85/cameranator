import "bootstrap/dist/css/bootstrap.min.css";
import React from "react";
import ReactDOM from "react-dom";
import App from "./components/app/app";
import { AppLogic } from "./components/app/app_logic";
import { AppProps } from "./components/app/app_props";
import { Alert, Container, Row, Spinner } from "react-bootstrap";

new AppLogic((props) => {
    ReactDOM.render(<App {...props} />, document.getElementById("root"));
});

interface IndexState {
    appProps: AppProps | null;
}

class Index extends React.Component<any, IndexState> {
    mounted: boolean;
    appModel: AppLogic;

    constructor(props: any, state: IndexState) {
        super(props, state);

        this.state = {
            appProps: null,
        };

        this.mounted = false;
        this.appModel = new AppLogic((appProps) => {
            this.appUpdateHandler(appProps);
        });
    }

    appUpdateHandler(appProps: AppProps) {
        if (!this.mounted) {
            return;
        }

        this.setState({ appProps: appProps });
    }

    componentDidMount() {
        this.mounted = true;
    }

    componentWillUnmount() {
        this.mounted = false;
    }

    render() {
        if (this.state.appProps == null) {
            return (
                <Container>
                    <Row className="justify-content-md-center">
                        <Alert variant={"info"}>Loading, please wait...</Alert>
                    </Row>
                    <Row className="justify-content-md-center">
                        <Spinner animation="border" role="status">
                            <span className="sr-only">Loading...</span>
                        </Spinner>
                    </Row>
                </Container>
            );
        }

        return <App {...this.state.appProps} />;
    }
}

ReactDOM.render(<Index />, document.getElementById("root"));
