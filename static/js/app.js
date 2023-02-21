function Prompt() {
    let toast = function (c) {
        const {
            title = "",
            icon = "success",
            position = "top-end",
        } = c;

        const Toast = Swal.mixin({
            toast: true,
            title: title,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        });

        Toast.fire({});
    }

    let success = function (c) {
        const {
            title = "",
            text = "",
            footer = ""
        } = c;
        Swal.fire({
            icon: 'success',
            title: title,
            text: text,
            footer: footer
        })
    }

    let error = function (c) {
        const {
            title = "",
            text = "",
            footer = ""
        } = c;
        Swal.fire({
            icon: 'error',
            title: title,
            text: text,
            footer: footer
        })
    }

    async function custom(c) {
        const {
            icon = "",
            msg = "",
            title = "",
            showConfirmButton = true,
        } = c;

        const { value: result } = await Swal.fire({
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            icon: icon,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            }
        })

        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== "") {
                    if (c.callback !== undefined) {
                        c.callback(result);
                    }
                } else {
                    c.callback(false);
                }
            } else {
                c.callback(false);
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}

function createReservationModal(room, csrf_token) {
    let html = `
        <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
                <div class="col">
                    <div class="form-row" id="reservation-dates-modal">
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="start_modal" id="start_modal" placeholder="Arrival">
                        </div>
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="end_modal" id="end_modal" placeholder="Departure">
                        </div>
                    </div>
                </div>
        </form>
            `

    attention.custom({
        msg: html,
        title: 'Choose your dates',
        willOpen: () => {
            const elem = document.getElementById('reservation-dates-modal');
            const rp = new DateRangePicker(elem, {
                format: 'yyyy-mm-dd',
                showOnFocus: true,
                minDate: new Date(),
            });
        },
        didOpen: () => {
            document.getElementById('start_modal').removeAttribute('disabled');
            document.getElementById('end_modal').removeAttribute('disabled');
        },
        callback: function (result) {
            let form = document.getElementById("check-availability-form");
            let formData = new FormData(form);
            formData.append("csrf_token", csrf_token);
            formData.append("room_id", room);
            console.log(formData);

            fetch('/search-availability-json', {
                method: "post",
                body: formData,
            })
                .then(response => response.json())
                .then(data => {
                    console.log(data)
                    if (data.ok) {
                        attention.custom({
                            icon: 'success',
                            msg: '<p>Room is available</p>' +
                                '<p><a href="/book-room?id=' + data.room_id + '&s=' + data.start_date + '&e=' + data.end_date +
                                '" class="btn btn-primary">Book now!</a></p>',
                            showConfirmButton: false,
                        })
                    } else {
                        attention.error({
                            title: "No availability"
                        });
                    }
                })
        }
    })
}