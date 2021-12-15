import "bootstrap/dist/css/bootstrap.min.css";
import React from "react";
import ReactDOM from "react-dom";
import App from "./components/app/app";
import { AppLogic } from "./components/app/app_logic";
import { AppProps } from "./components/app/app_props";
import { Alert, Container, Row, Spinner } from "react-bootstrap";
import { info } from "./common/utils";

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
        info(`${this.constructor.name}.appUpdateHandler fired`);

        if (!this.mounted) {
            return;
        }

        this.setState({ appProps: appProps });
    }

    async componentDidMount() {
        info(`${this.constructor.name} mounted`);

        this.mounted = true;
        await this.appModel.updateAll();
    }

    componentWillUnmount() {
        info(`${this.constructor.name} unmounted`);

        this.mounted = false;
    }

    render() {
        info(`${this.constructor.name}.appUpdateHandler rendering`);

        if (this.state.appProps == null) {
            info(
                `${this.constructor.name}.appUpdateHandler showing loading banner`
            );
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

        info(`${this.constructor.name}.appUpdateHandler showing App`);
        return <App {...this.state.appProps} />;
    }
}

ReactDOM.render(<Index />, document.getElementById("root"));
