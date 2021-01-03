import moment from "moment/moment";
import { Camera, getCameras } from "../../persistence/types/camera";
import { Event, getEvents } from "../../persistence/types/event";
import { getDates } from "../../persistence/types/date";
import { SEGMENTS } from "../type_drop_down/type_drop_down";
import { AppProps } from "./app_props";

type AppLogicUpdateHandler = {
    (props: AppProps): void;
};

export class AppLogic {
    connected: boolean;
    dates: moment.Moment[];
    cameras: Camera[];
    events: Event[];
    type: string | null;
    date: moment.Moment | null;
    camera: Camera | null;
    handler: AppLogicUpdateHandler;

    constructor(handler: AppLogicUpdateHandler) {
        this.connected = false;
        this.dates = [] as moment.Moment[];
        this.cameras = [] as Camera[];
        this.events = [] as Event[];
        this.type = null;
        this.date = null;
        this.camera = null;

        this.handler = handler;

        setInterval(() => {
            this.update();
        }, 1000);
    }

    updateDates() {
        getDates((dates) => {
            if (dates === null) {
                this.connected = false;
            } else {
                this.dates = dates;
                this.connected = true;
            }

            this.handler(this.getAppProps());
        });
    }

    updateCameras() {
        getCameras((cameras) => {
            if (cameras === null) {
                this.connected = false;
            } else {
                this.cameras = cameras;
                this.connected = true;
            }

            this.handler(this.getAppProps());
        });
    }

    updateEvents() {
        if (!(this.connected && this.type && this.date && this.camera)) {
            return;
        }

        getEvents(
            this.type === SEGMENTS,
            this.date,
            this.camera.name,
            (events) => {
                if (events === null) {
                    this.connected = false;
                } else {
                    this.events = events;
                    this.connected = true;
                }

                this.handler(this.getAppProps());
            }
        );
    }

    update() {
        this.updateDates();

        this.updateCameras();

        this.updateEvents();
    }

    public setType(type: string) {
        this.type = type;
        this.update();
    }

    public setDate(date: moment.Moment) {
        this.date = date;
        this.update();
    }

    public setCamera(camera: Camera) {
        this.camera = camera;
        this.update();
    }

    getAppProps(): AppProps {
        return {
            connected: this.connected,
            dates: this.dates,
            cameras: this.cameras,
            events: this.events,
            type: this.type,
            date: this.date,
            camera: this.camera,
            typeChangeHandler: (type: any) => {
                this.setType(type);
            },
            dateChangeHandler: (date: any) => {
                this.setDate(date);
            },
            cameraChangeHandler: (camera: any) => {
                this.setCamera(camera);
            },
        };
    }
}
