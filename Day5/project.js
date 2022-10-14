let dataProject= new Array;

function addProject(){
    let namaProject = document.getElementById("project-name").value;
    let startDate = document.getElementById("start-date").value;
    let endDate = document.getElementById("end-date").value;

    let date1 = startDate.split('-');
    let date2 = endDate.split('-');

    let setdate1 = new Date(date1[0], date1[1], date1[2]);
    let setdate2 = new Date(date2[0], date2[1], date2[2]);

    let time1 = parseInt(setdate1.getTime()/1000);
    let time2 = parseInt(setdate2.getTime()/1000);

    let detik = time2 - time1;
    let jam = Math.round(detik/60/60);
    let hari = Math.round(jam/24);
    let bulan = Math.round(hari/30);
    let tahun = Math.round(bulan/12);

    console.log(detik)
    console.log(jam)
    console.log(hari)
    console.log(bulan)
    console.log(tahun)
    let durasi; 

    if(jam < 24){
        durasi = "Hari ini "; 
    }else if(jam>=24 && jam < 720){
        durasi = hari+" Hari";
    }else if(jam >= 720 && jam < 8760){
        durasi = bulan+" Bulan";
    }else{
        durasi= tahun+" Tahun";
    }

    let deskripsi = document.getElementById("description").value;
    let viewdeskripsi;
    
    if(deskripsi.length > 150){
        viewdeskripsi = deskripsi.substring(0,150)+"....";
    }else{
        viewdeskripsi = deskripsi;
    }
     
    
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
        viewdeskripsi,
        image,
        tehnik,
        waktupost : new Date(),
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
        <div class="blog-container">
            <a href="blog.html" target="blank">
                <img src="${dataProject[index].image}">
                <h3>${dataProject[index].namaProject}</h3>
                <label>Durasi : ${dataProject[index].durasi}</label>
                <p>${dataProject[index].viewdeskripsi}</p>
                <div class="toggle">
                    <a href="#"><i class="fa-brands fa-react"></i></i></i></a>
                    <a href="#"><i class="fa-brands fa-android"></i></a>
                    <a href="#"><i class="fa-brands fa-square-js"></i></a>
                    <label>Post : ${timePostProject(dataProject[index].waktupost)}</label>
                </div>
                <div style="display: flex;justify-content:center ;">
                    <button>Edit</button>
                    <button>Hapus</button>
                </div>
            </a>
        </div>
        `
    }
}
function timePostProject(timeproject) {
    let timeNow = new Date();
    let timePost = timeproject;

    let time = timeNow - timePost;
    
    let satuanMiliDetik = 1000;
    let satuanDetikDalamMenit = 60;
    let satuanDetikDalamJam = 3600 ;
    let satuanJamDalamSehari = 24

    let hitungHari = Math.floor(time / (satuanMiliDetik * satuanDetikDalamJam * satuanJamDalamSehari))
    let hitungJam = Math.floor(time / (satuanMiliDetik * satuanDetikDalamJam))
    let hitungMenit = Math.floor(time / (satuanMiliDetik * satuanDetikDalamMenit))
    let hitungDetik = Math.floor(time / satuanMiliDetik)

    if (hitungHari > 0) {
        return `${hitungHari} Hari yang lalu`
    } else if (hitungJam > 0) {
        return `${hitungJam} Jam yang lalu`
    } else if (hitungMenit > 0) {
        return `${hitungMenit} Menit yang lalu`
    } else {
        return `Baru saja`
    }
}
setInterval(function() {
    renderProject()
}, 60000);//hitung per menit

