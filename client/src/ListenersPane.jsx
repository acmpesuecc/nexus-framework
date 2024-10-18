import { useEffect, useState } from 'react';
import axios from 'axios';

const ListenersPane = () => {
    const [listeners, setListeners] = useState([]);
    const [lhost, setLhost] = useState('');
    const [lport, setLport] = useState('');

    const fetchListeners = async () => {
        const response = await axios.get('/api/listeners');
        setListeners(response.data.listeners);
    };

    const createListener = async (e) => {
        e.preventDefault(); // Prevent the default form submission
        await axios.post('/api/listeners', { lhost, lport });
        fetchListeners(); // Refresh the list after creating a new listener
    };

    useEffect(() => {
        fetchListeners();
    }, []);

    return (
        <div>
            <h1>Listeners</h1>
            {listeners.map((listener, index) => (
                <div key={index}>{`${listener.host}:${listener.port}`}</div>
            ))}
            <form onSubmit={createListener}>
                <input
                    type="text"
                    placeholder="Host"
                    value={lhost}
                    onChange={(e) => setLhost(e.target.value)}
                />
                <input
                    type="number"
                    placeholder="Port"
                    value={lport}
                    onChange={(e) => setLport(e.target.value)}
                />
                <button type="submit">Add Listener</button>
            </form>
        </div>
    );
};

export default ListenersPane;
