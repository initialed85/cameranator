import React from "react";
import { NavDropdown } from "react-bootstrap";
import moment from "moment/moment";

export interface DateDropDownChangeHandler {
    (date: moment.Moment): void;
}

interface DateDropdownProps {
    dates: moment.Moment[];
    changeHandler: DateDropDownChangeHandler;
}

export class DateDropdown extends React.Component<DateDropdownProps, any> {
    render() {
        let items: any[] = [];

        this.props.dates.forEach((date: moment.Moment) => {
            const friendlyDate = date.format("YYYY-MM-DD");

            items.push(
                <NavDropdown.Item
                    style={{ fontSize: "10pt" }}
                    href={`#date/${friendlyDate}`}
                    key={friendlyDate}
                    onClick={() => {
                        this.props.changeHandler(date);
                    }}
                >
                    {friendlyDate}
                </NavDropdown.Item>
            );
        });

        return (
            <NavDropdown title="Date" id="date-dropdown">
                {items}
            </NavDropdown>
        );
    }
}
