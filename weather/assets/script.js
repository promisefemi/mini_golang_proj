console.log("adjhshj")

let button = document.querySelector("#changeCityButton")
let modal = document.querySelector(".modal")
let modalBackdrop = document.querySelector(".modalBackdrop")
let results = document.querySelector(".results")
let input = document.querySelector(".modal input")


let cities = document.querySelector(".modal textarea").value
cities = JSON.parse(cities)
console.log(cities)

button.addEventListener("click", function () {
    toggleModal()

})

modalBackdrop.addEventListener("click", function () {
    toggleModal()

})



input.addEventListener('keyup', function (e) {
    filteredCities = []

    if (e.target.value.trim() != "") {
        filteredCities = cities.filter(function (city) {
            return city.Name.toLowerCase().match(e.target.value.toLowerCase().trim())

        })
    }


    renderResults(filteredCities)
})


function renderResults(cities) {

    let resultsContent = ''
    cities.forEach(function (city) {
        resultsContent += `<a href="/?cityid=${city.ID}">${city.Name}</a>`
    })


    results.innerHTML = resultsContent

}




function toggleModal() {

    if (modal.classList.contains("showModal")) {

        modal.classList.remove("showModal")
        modalBackdrop.classList.remove("showModal")

    } else {
        modal.classList.add("showModal")
        modalBackdrop.classList.add("showModal")
    }
}