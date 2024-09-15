async function fetchData() {
    return fetch("/api/data")
        .then(response => response.json())
        .then(data => {
            return data;
        })
        .catch(error => {
            console.error("Error fetching data:", error);
            throw error;
        });
}

function getRandomIndexForLocalStorage(length) {
    const randomIndex = Math.floor(Math.random() * length);
    return randomIndex;
}

function popIDFromLocalStorage() {
    const keys = Object.keys(localStorage);
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

document.addEventListener("DOMContentLoaded", function() {
    const generateButton = document.getElementById("generate-button");
    generateButton.addEventListener("click", function() {
        const randomID = popIDFromLocalStorage();
        console.log(randomID);
    });

    if (localStorage.length === 0) {
        fetchData()
            .then(data => {
                if (data) {
                    console.log(data);
                    setData(data.ids);
                }
            })
            .catch(error => {
                console.error("Error setting data:", error);
            });
    }
})
