var navigator = document.querySelectorAll("nav a");
for (var i = 0; i < navigator.length; i++) {
	var j = navigator[i]
	if (j.getAttribute('href') == window.location.pathname) {
		j.classList.add("live");
		break;
	}
}

let comment_incr = document.querySelector(".add__comment")
let comm=document.querySelector(".comment")
comment_incr.addEventListener("click", ()=> {
	comm.classList.toggle("active")
})