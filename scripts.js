// var jwt = localStorage.getItem("jwt");

//       if (jwt == null) {
//         window.location.href = 'URL';
//       }

$(document).ready(function () {
    var username = 'Username';
    $('#username').text(username);

    $('#score-form').submit(function (e) {
        e.preventDefault();

        var score = $('#score').val();
        var name = $('#name').val();
        var data = JSON.stringify({
            username: username,
            score: parseInt(score),
            name: name
        });

        $.ajax({
            type: 'POST',
            url: 'http://localhost:8080/grade',
            data: data,
            contentType: 'application/json',
            success: function (response) {
                console.log(response); // ตรวจสอบ response ใน console.log เพื่อดูค่าที่ได้รับ
                if (response.status === 'OK') {
                    Swal.fire({
                        position: "top-end",
                        icon: "success",
                        title: "Your work has been saved",
                        showConfirmButton: false,
                        timer: 1500
                    });
                } else {
                    Swal.fire({
                        icon: "error",
                        title: "",
                        text: "Failed to save data",
                    });
                }
            },
            error: function (xhr, status, error) {
                Swal.fire({
                    icon: "error",
                    title: "",
                    text: "No connection!",
                });
            }
        });
    });
});
