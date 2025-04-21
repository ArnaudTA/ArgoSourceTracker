<script setup lang="ts">
import { Button, Panel } from 'primevue';

</script>

<template>

    <div id="page">
        <Panel header="Swagger">
                <Button as="a" href="/swagger/index.html">EXPLORE API</Button>
        </Panel>
        <div class="element">
            <p>Source Code</p>
            <Button as="a" href="https://github.com/cableship/argocd-source-tracker">EXPLORE CODE</Button>
        </div>
        <div class="element">
            <p>ArtifactHub</p>
            <Button as="a" href="https://artifacthub.io/packages/helm/cableship/argocd-source-tracker">VIEW ON
                ARTIFACTHUB</Button>
        </div>
    </div>
</template>
<style scoped>
#page {
    margin: 3rem;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 300px));

    flex-direction: row;
    flex-basis: 25%;
    height: min-content;
    gap: 3rem;
    flex-wrap: wrap;
    justify-content: space-evenly;
}

.element {
    display: flex;
    flex-direction: column;
    align-content: center;
    padding: 2rem;
    gap: 1rem;
    /* background-color: white; */
    text-align: center;
    box-shadow: 0 0 0.3rem 0.3rem rgba(128, 128, 128, 0.356);
}

.p-button {
    padding: 1rem;
    /* background: #6d7f8b; */
    border-radius: 1.5rem;
    font-weight: bolder;
    text-decoration: none;
}
</style>