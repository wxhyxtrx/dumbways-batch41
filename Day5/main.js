function sendMail(){
    let name = document.getElementById('input-name').value;
    let mail = document.getElementById('input-email').value;
    let phone = document.getElementById('input-phone').value;
    let subject = document.getElementById('input-subject').value;
    let message = document.getElementById('input-message').value;

    console.log(name , phone,mail , subject, message)

    if(name =='' || mail=='' || phone == '' || subject == '' || message== ''){
        return alert("Tolong isi dengan lengkap!");
    }else{
        let penerimapesan = "banyuwangi1999@gmail.com";

        let inbox = document.createElement('a');
        inbox.href =`https://mail.google.com/mail/?view=cm&fs=1&to=${penerimapesan}&su=${subject}&body=hallo nama saya ${name} ${message} silahkan hubungi ${phone}`;

        inbox.click()
        
    }
}