let dataProject= new Array;

function addProject(){
    let namaProject = document.getElementById("project-name").value;
    let startDate = document.getElementById("start-date").value;
    let endDate = document.getElementById("end-date").value;

    let date1 = startDate.split('-');
    let date2 = endDate.split('-');

    let year = (date2[0] - date1[0]);
    let month = (date2[1] - date1[1]);
    let day = (date2[2] - date1[2]);

    let setdate1 = new Date(date1[0], date1[1], date1[2]);
    let setdate2 = new Date(date2[0], date2[1], date2[2]);

    let time1 = parseInt(setdate1.getTime()/1000);
    let time2 = parseInt(setdate2.getTime()/1000);

    let detik = time2 - time1;
    let jam = Math.round(detik/60/60);
    let hari = Math.round(jam/24);
    let bulan = Math.round(hari/30);
    let tahun = Math.round(bulan/12);

    let durasi; 

    if(jam <= 24){
        durasi = jam+" Jam"; 
    }else if(hari<=30){
        durasi = hari+" Hari";
    }else if(hari>30){
        durasi = bulan+" Bulan";
    }else if(bulan > 12){
        durasi= tahun+" Tahun";
    }

    let deskripsi = document.getElementById("description").value;
    
    let node = document.getElementById("node").value;
    let typescript = document.getElementById("typescript").value;
    let react = document.getElementById("react").value;
    let next = document.getElementById("next").value;

    let tehnik = node+" "+typescript+" "+react+" "+next;
        
    let filegambar = document.getElementById("input-image").files[0];
    let image = URL.createObjectURL(filegambar)

    let project = {
        namaProject,
        durasi,
        deskripsi,
        image,
        tehnik,
        author: "Wahyu Tricahyo"
    }

    dataProject.push(project);
    renderProject();

    alert("Berhasil Di Tambahkan");
}
function renderProject() {
    document.getElementById("conten-project").style.backgroundColor = "white";
    document.getElementById("view").innerHTML = ""

    for (let index = 0; index < dataProject.length; index++) {
        document.getElementById("view").innerHTML += `
    <a href="blog.html" target="blank">
    <div class="blog-container">
        <img src="${dataProject[index].image}">
        <h3>${dataProject[index].namaProject}</h3>
        <label>Durasi : ${dataProject[index].durasi}</label>
        <p>${dataProject[index].deskripsi} </p>
        <div class="toggle">
            <a href="#"><i class="fa-brands fa-react"></i></i></i></a>
            <a href="#"><i class="fa-brands fa-android"></i></a>
            <a href="#"><i class="fa-brands fa-square-js"></i></a>
        </div>
        <div style="display: flex;justify-content:center ;">
            <button>Edit</button>
            <button>Hapus</button>
        </div>
    </div>
    </a>
        `
    }
}

