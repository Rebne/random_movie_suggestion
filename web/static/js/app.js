function getRandomIndexForLocalStorage(length) {
    const randomIndex = Math.floor(Math.random() * length);
    return randomIndex;
}

function popIDFromLocalStorage() {
    const keys = Object.keys(localStorage).filter(key => key !== 'totalIds');
    if (keys.length === 0) {
        return null;
    }
    const randomIndex = getRandomIndexForLocalStorage(keys.length);
    const id = keys[randomIndex];
    localStorage.removeItem(id);
    return id;
}

function setData(ids) {
    if (ids && Array.isArray(ids)) {
        ids.forEach(id => {
            localStorage.setItem(id, id);
        });
    }
}

async function initializeLocalStorage() {
    try {
        if (localStorage.length <= 1) {
            const response = await fetch('/api/data');
            const data = await response.json();
            const ids = data.ids;
            if (localStorage.length == 0) {
                localStorage.setItem('totalIds', data.total.toString());
            }
            setData(ids);
        } else {
            const response = await fetch('/api/data/length');
            const data = await response.json();
            const current = parseInt(localStorage.getItem('totalIds'));
            if (current < data.length) {
                const updateResponse = await fetch('/api/data/new', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ currentLength: localStorage.getItem('totalIds') })
                });
                const updateData = await updateResponse.json();
                setData(updateData.newIDs);
                localStorage.setItem('totalIds', updateData.newLength);
            }
        }
        document.body.dispatchEvent(new Event('local-storage-updated'));
    } catch (error) {
        console.error('Error occurred fetching API data from server: ', error);
    }
}