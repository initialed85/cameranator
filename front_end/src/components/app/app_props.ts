import moment from "moment/moment";
import { Camera } from "../../persistence/collections/camera";
import { Event } from "../../persistence/collections/event";

export interface AppProps {
    connected: boolean;
    dates: moment.Moment[];
    cameras: Camera[];
    events: Event[];
    type: string | null;
    date: moment.Moment | null;
    camera: Camera | null;
    typeChangeHandler: CallableFunction;
    dateChangeHandler: CallableFunction;
    cameraChangeHandler: CallableFunction;
}
