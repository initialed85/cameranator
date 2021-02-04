import moment from "moment/moment";
import {SEGMENTS} from "../type_drop_down/type_drop_down";
import {AppProps} from "./app_props";
import DateCollection from "../../persistence/collections/date";
import CameraCollection, {Camera} from "../../persistence/collections/camera";
import {Event, EventCollection} from "../../persistence/collections/event";

type AppLogicUpdateHandler = {
    (props: AppProps): void;
};

export class AppLogic {
    dateCollection: DateCollection;
    cameraCollection: CameraCollection;
    eventCollection: EventCollection;
    connected: boolean;
    dates: moment.Moment[];
    cameras: Camera[];
    events: Event[];
    type: string | null;
    date: moment.Moment | null;
    camera: Camera | null;
    handler: AppLogicUpdateHandler;

    constructor(handler: AppLogicUpdateHandler) {
        this.dateCollection = new DateCollection();
        this.cameraCollection = new CameraCollection();
        this.eventCollection = new EventCollection();
        this.connected = false;
        this.dates = [] as moment.Moment[];
        this.cameras = [] as Camera[];
        this.events = [] as Event[];
        this.type = null;
        this.date = null;
        this.camera = null;

        this.handler = handler;

        setInterval(() => {
            if (!(this.type && this.date && this.camera)) {
                return;
            }

            this.eventCollection.get({
                isSegment: this.type === SEGMENTS,
                date: this.date,
                cameraName: this.camera.name,
            }).then((events) => {
                if (!events) {
                    return;
                }

                this.events = events;
                this.handler(this.getAppProps());
            })
        }, 10000);
    }

    public updateAll() {
        this.cameraCollection
            .get({})
            .catch((e) => {
                this.connected = false;
                console.error(e);
            })
            .then((cameras) => {
                if (!cameras) {
                    return;
                }

                this.connected = true;
                this.cameras = cameras;

                return this.dateCollection.get({});
            })
            .then((dates) => {
                if (!dates) {
                    return;
                }

                this.dates = dates;

                if (!(this.type && this.camera && this.date)) {
                    return;
                }

                return this.eventCollection.get({
                    isSegment: this.type === SEGMENTS,
                    date: this.date,
                    cameraName: this.camera.name,
                });
            })
            .then((events) => {
                if (!events) {
                    return;
                }

                this.events = events;
            })
            .then(() => {
                this.handler(this.getAppProps());
            });
    }

    public setType(type: string | null) {
        if (type === null) {
            this.events = [];
        }

        this.type = type;
        this.updateAll();
    }

    public setDate(date: moment.Moment | null) {
        if (date === null) {
            this.events = [];
        }

        this.date = date;
        this.updateAll();
    }

    public setCamera(camera: Camera | null) {
        if (camera === null) {
            this.events = [];
        }

        this.camera = camera;
        this.updateAll();
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
