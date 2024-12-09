<script lang="ts">
    export let issueKey: string;
    export let baseUrl: string;
    import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
    import { Button } from "flowbite-svelte";
</script>

<Button
    class="inline-block px-2 py-1 bg-gray-100 text-blue-700 font-semibold rounded-md hover:bg-gray-200 transition"
    on:click={() => BrowserOpenURL(`${baseUrl}/browse/${issueKey}`)}
>
    <span class="line-clamp-1">
        {issueKey}
    </span>
</Button>
