import moment from "moment/moment";
import { SEGMENTS } from "../type_drop_down/type_drop_down";
import { AppProps } from "./app_props";
import DateCollection from "../../persistence/collections/date";
import CameraCollection, { Camera } from "../../persistence/collections/camera";
import { Event, EventCollection } from "../../persistence/collections/event";
import { error, info, warn } from "../../common/utils";

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
    handler: () => void;

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

        this.handler = () => {
            handler(this.getAppProps());
        };

        setInterval(async () => {
            info(`${this.constructor.name}.setInterval fired`);
            await this.updateEvents();
        }, 10000);
    }

    private async updateCamerasAndDatesAndConnected() {
        info(
            `${this.constructor.name}.updateCamerasAndDatesAndConnected fired`
        );

        const camerasPromise = this.cameraCollection.get({});
        const datesPromise = this.dateCollection.get({});

        let cameras;
        let dates;

        try {
            cameras = await camerasPromise;
            dates = await datesPromise;
        } catch (e) {
            this.connected = false;
            error(
                `${this.constructor.name}.updateCamerasAndDatesAndConnected ${e}`
            );
            return;
        }

        this.connected = true;

        info(
            `${
                this.constructor.name
            }.updateCamerasAndDatesAndConnected; cameras = ${
                cameras?.length || 0
            }`
        );
        if (cameras?.length) {
            this.cameras = cameras;
        }

        info(
            `${
                this.constructor.name
            }.updateCamerasAndDatesAndConnected; dates = ${dates?.length || 0}`
        );
        if (dates?.length) {
            this.dates = dates;
        }
    }

    private async updateEvents() {
        info(`${this.constructor.name}.updateEvents fired`);

        if (!(this.cameras?.length && this.dates?.length)) {
            warn(
                `${this.constructor.name}.updateCamerasAndDatesAndConnected; cameras and / or dates missing`
            );
            return;
        }

        if (!(this.type && this.camera && this.date)) {
            warn(
                `${this.constructor.name}.updateCamerasAndDatesAndConnected; type and / or camera and / or date not selected`
            );
            return;
        }

        let events;

        try {
            events = await this.eventCollection.get({
                isSegment: this.type === SEGMENTS,
                date: this.date,
                cameraName: this.camera.name,
            });
        } catch (e) {
            error(`${this.constructor.name}.updateEvents ${e}`);
            return;
        }

        info(
            `${this.constructor.name}.updateEvents; events = ${
                events?.length || 0
            }`
        );
        if (!events?.length) {
            return;
        }

        this.events = events;
    }

    public async updateAll() {
        info(`${this.constructor.name}.updateAll fired`);
        await this.updateCamerasAndDatesAndConnected();
        await this.updateEvents();
        this.handler();
    }

    public async setType(type: string | null) {
        info(`${this.constructor.name}.setType fired`);
        if (type === null) {
            this.events = [];
        }

        this.type = type;
        await this.updateAll();
    }

    public async setDate(date: moment.Moment | null) {
        info(`${this.constructor.name}.setDate fired`);
        if (date === null) {
            this.events = [];
        }

        this.date = date;
        await this.updateAll();
    }

    public async setCamera(camera: Camera | null) {
        info(`${this.constructor.name}.setCamera fired`);
        if (camera === null) {
            this.events = [];
        }

        this.camera = camera;
        await this.updateAll();
    }

    getAppProps(): AppProps {
        info(`${this.constructor.name}.getAppProps fired`);
        return {
            connected: this.connected,
            dates: this.dates,
            cameras: this.cameras,
            events: this.events,
            type: this.type,
            date: this.date,
            camera: this.camera,
            typeChangeHandler: async (type: any) => {
                await this.setType(type);
            },
            dateChangeHandler: async (date: any) => {
                await this.setDate(date);
            },
            cameraChangeHandler: async (camera: any) => {
                await this.setCamera(camera);
            },
        };
    }
}
