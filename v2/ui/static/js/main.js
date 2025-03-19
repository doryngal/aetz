var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

const selectBtn = document.querySelector(".select-btn");
const items = document.querySelectorAll(".item");
// const btnTexts = document.querySelectorAll(".btn-text");



function  SelectedPlace() {
    let checked = document.querySelectorAll(".checked");
    btnText = document.querySelector('.btn-text')

    if (checked && checked.length > 0 ) {
        btnText.innerText = `Выбрано: ${checked.length}`
    }else {
        btnText.innerText = "Выберите место"
    }
}


selectBtn.addEventListener("click", () => {
    selectBtn.classList.toggle("open");
});


items.forEach(item => {
    const checkbox = item.querySelector('.checkbox')


    // Обработчик клика по checkbox, предотвращаем его перекрытие
    checkbox.addEventListener("click", (event) => {
        event.stopPropagation(); // Останавливаем всплытие события, чтобы клик не обрабатывался дальше
        item.classList.toggle("checked", checkbox.checked); // Обновляем класс checked в зависимости от состояния чекбокса
        SelectedPlace(); // Обновляем текст
    });

    item.addEventListener("click", () => {
        checkbox.checked = !checkbox.checked; // Переключаем состояние checkbox
        item.classList.toggle("checked"); // Устанавливаем класс на item
        SelectedPlace(); // Обновляем текст
    });
});



// Обработка кнопки "Выбрать все"
const checkAll = document.getElementById("check_all");
checkAll.addEventListener("click", () => {
    items.forEach(item => {
        const checkbox = item.querySelector('.checkbox');
        checkbox.checked = true; // Устанавливаем checked на true
        item.classList.add("checked"); // Устанавливаем класс на item
    });
    SelectedPlace(); // Обновляем текст
});

// Обработка кнопки "Сбросить"
const reset = document.getElementById("reset");
reset.addEventListener("click", () => {
    items.forEach(item => {
        const checkbox = item.querySelector('.checkbox');
        checkbox.checked = false; // Устанавливаем checked на false
        item.classList.remove("checked"); // Убираем класс на item
    });
    SelectedPlace(); // Обновляем текст
});



const radios = document.querySelectorAll('.radio-lot');
const submitButton = document.getElementById('review_confirm_lot');
const cancelReview = document.getElementById('cancel_review') 

// Добавляем обработчик события на изменение радио-кнопок
radios.forEach(radio => {
    radio.addEventListener('change', function() {
        // Если любая радио-кнопка выбрана, разблокируем кнопку
        submitButton.disabled = false;
        cancelReview.disabled = false
    });
});

cancelReview.addEventListener("click",() => {
    radios.forEach(radio => {
        radio.checked = false; // Отменяем выбор радиокнопки
        submitButton.disabled = true;
        cancelReview.disabled = true;
    });
})


const back = document.getElementById('backButton') 

back.addEventListener('click', function() {
    console.log(1)
    window.history.back();
});





// const selectBtn = document.querySelector(".select-btn")
// items = document.querySelectorAll(".item")

// selectBtn.addEventListener("click",()=>{
// 	selectBtn.classList.toggle("open")
// })

// items.forEach(item => {
// 	item.addEventListener("click",() =>{
// 		item.classList.toggle("checked")

// 		let checked  = document.querySelectorAll(".checked");
// 		btnText = document.querySelectorAll(".btn-text");
// 		let selectedTexts = Array.from(checked).map(item => {
//             return item.querySelector(".item-text").innerText;
//         });
// 		console.log(selectedTexts)
// 		if (checked.length > 0) {
// 			btnText.forEach(btnText => {
// 				btnText.innerText = `Выбрано место поставок ${selectedTexts.join(', ')}`;
// 			});
// 		} else {
// 			btnText.forEach(btnText => {
// 				btnText.innerText = 'Выбрано место поставок'; 
// 			});
// 		}
// 	})
// })