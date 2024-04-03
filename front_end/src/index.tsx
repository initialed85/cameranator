import React from "react"
import "./index.css"
import reportWebVitals from "./reportWebVitals"
import { ApolloProvider } from "@apollo/client"
import { client } from "./hasura/client"
import { createRoot } from "react-dom/client"
import App from "./components/app/App"

const addMaximumScaleToMetaViewport = () => {
    const el = document.querySelector("meta[name=viewport]")

    if (el !== null) {
        let content = el.getAttribute("content")
        if (content !== null) {
            let re = /maximum\-scale=[0-9\.]+/g

            if (re.test(content)) {
                content = content.replace(re, "maximum-scale=1.0")
            } else {
                content = [content, "maximum-scale=1.0"].join(", ")
            }

            el.setAttribute("content", content)
        }
    }
}

const disableIosTextFieldZoom = addMaximumScaleToMetaViewport

// https://stackoverflow.com/questions/9038625/detect-if-device-is-ios/9039885#9039885
const checkIsIOS = () =>
    /iPad|iPhone|iPod/.test(navigator.userAgent) && !(window as any).MSStream

if (checkIsIOS()) {
    disableIosTextFieldZoom()
}

const root = createRoot(document.getElementById("root")!)

root.render(
    <React.StrictMode>
        <ApolloProvider client={client}>
            <App />
        </ApolloProvider>
    </React.StrictMode>,
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()
