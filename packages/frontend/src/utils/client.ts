import { Api } from "../api/Api";
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