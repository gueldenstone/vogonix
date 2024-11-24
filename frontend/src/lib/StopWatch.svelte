<script lang="ts">
    import { Button } from "flowbite-svelte";
    import * as Icons from "flowbite-svelte-icons";
    import { onMount } from "svelte";
    import { formatDuration } from "@lib/util";

    // export let time: number = 0;
    export let isRunning: boolean = false;

    export let startCallback;
    export let pauseCallback;
    export let resetCallback;
    export let setupCallback;
    // export let newTimeCallback;

    export let time: number = 0;

    onMount(() => {
        setupCallback();
    });
</script>

<div class="stopwatch">
    <!-- Display the timer -->
    <div>
        {formatDuration(time)}
    </div>

    <!-- Control buttons -->
    <div class="buttons">
        {#if !isRunning}
            <Button
                on:click={() => {
                    startCallback();
                    isRunning = true;
                }}
                color="light"
                size="xs"
            >
                <Icons.PlaySolid />
            </Button>
        {:else}
            <Button
                on:click={() => {
                    pauseCallback();
                    isRunning = false;
                }}
                color="light"
                size="xs"
            >
                <Icons.PauseSolid />
            </Button>
        {/if}
        <Button
            on:click={() => {
                resetCallback();
                isRunning = false;
            }}
            color="light"
            size="xs"
        >
            <Icons.CloseCircleSolid />
        </Button>
    </div>
</div>
