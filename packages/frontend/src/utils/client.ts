import { Api, TypesApplicationStatus } from "../api/Api";
import router from "../router";

export const client = new Api({
    baseURL: '/'
})

export function goToApp(application: { name: string, namespace: string }) {
    router.push({
        name: 'Application',
        params: {
            name: application.name,
            namespace: application.namespace
        }
    })
}

export type TileStatus = 'Up-to-date' | 'Ignored' | 'Outdated' | 'Checking' | 'Unknown' | 'Error' | 'None'

export const statusClass: Record<TypesApplicationStatus | 'None', string> = {
    "Up-to-date": "uptodate",
    Ignored: "ignored",
    Outdated: "outdated",
    None: "none"
}