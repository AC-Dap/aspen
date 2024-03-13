import {createSignal, type Component, For} from 'solid-js';

const App: Component = () => {
    const [replies, setReplies] = createSignal<string[]>([]);

    const sendAPI = async () => {
        const response = await fetch('http://localhost:8080/api/test');
        const data = await response.text();
        setReplies((currReplies) => currReplies.concat(data));
    }

    return (
        <>
            <button onClick={sendAPI}>Send API</button>
            <For each={replies()}>
                {(reply) => <p>{reply}</p>}
            </For>
        </>
    );
};

export default App;
