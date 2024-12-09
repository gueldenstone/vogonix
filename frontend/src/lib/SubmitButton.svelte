<script lang="ts">
    import { Button } from "flowbite-svelte";
    import * as wails from "@runtime/runtime";
    import * as jira from "@go/jira/JiraInstance";
    import * as Icons from "flowbite-svelte-icons";

    export let enabled: boolean;
    export let issueKey: string;
</script>

{#if enabled}
    <Button
        size="xs"
        class="h-10"
        on:click={() =>
            jira.SubmitWorklog(issueKey).catch((err) => wails.LogError(err))}
        ><Icons.ShareAllSolid />
    </Button>
{:else}
    <Button disabled size="xs" class="h-10"><Icons.ShareAllSolid /></Button>
{/if}
