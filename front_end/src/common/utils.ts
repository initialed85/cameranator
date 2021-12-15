export function info(message: string): void {
    console.log(`${new Date().toISOString()} - ${message}`);
}

export function warn(message: string): void {
    console.warn(`${new Date().toISOString()} - ${message}`);
}

export function error(message: string): void {
    console.error(`${new Date().toISOString()} - ${message}`);
}
