const serverURL = "routes-ehgqedhpa7h5cpff.westeurope-01.azurewebsites.net"

$(function() {
    $("#region-combobox").combobox();
    $("#stop-combobox").combobox();

    locateUser();
});

function syncCombobox($select) {
    const $input = $select.next(".custom-combobox").find("input");
    const text = $select.find("option:selected").text();

    $input.val(text || "");
    if ($input.autocomplete("instance")) {
        $input.autocomplete("instance").term = text || "";
    }
}

// automatic prefill for region options
$(function() {
    $.get(`${serverURL}/regions`, function(data) {
        const $combobox = $("#region-combobox");
        
        $combobox.empty().append(
            $('<option>').val(" ").text("Select region...")
        );

        data.forEach(function(region) {
            $combobox.append(
                $('<option>').val(region).text(region)
            );
        });

        syncCombobox($combobox);
    })
})

function fillStops(region) {
    const $combobox = $("#stop-combobox");

    const url = `${serverURL}/stops?region=${encodeURIComponent(region)}`;

    return $.get(url, function(data) {
        $combobox.empty().append(
            $('<option>').val(" ").text("Select bus stop...")
        );

        data.forEach(function(stop) {
            $combobox.append(
                $('<option>')
                .val(stop.id)
                .data("stopName", stop.name)
                .data("stopCode", stop.code)
                .text(`${stop.name}, ${stop.code}`)
            );
        });
        syncCombobox($combobox);
        }).fail(function() {
        alert("Failed to fetch bus stops for region " + region);
    });
}


// get stops in regions on button click
$(document).on("click", "#region-combobox-button", function() {
    // empty and drop selected value of a bus stop container
    let $combobox = $("#stop-combobox");
    $combobox.empty().val(null);
    $combobox.next(".custom-combobox")
        .find("input")
        .val(null)

    // empty a routes container
    $combobox = $("#route-buttons-container");
    $combobox.empty();

    // empty a arrivals container
    $combobox = $("#arrivals-container");
    $combobox.empty();

    const region = $("#region-combobox").val();
    if (!region) {
      alert("Please select a region first!");
      return;
    }

    fillStops(region);
});

function fillRoutes(data) {
    const $container = $("#route-buttons-container");
    $container.empty(); // clear container from previous results
        
    data.forEach(function(route) {
        $container.append(
            $('<button>').addClass('route-button btn btn-outline-primary ').val(route.short_name).text(`${route.short_name}`)
        );
    });
}

// get routes for specific stop on button click
$(document).on("click", "#stop-combobox-button", function() {
    // empty a arrivals container
    const $combobox = $("#arrivals-container");
    $combobox.empty();

    const stop = $("#stop-combobox option:selected").data("stopName");
    const code = $("#stop-combobox option:selected").data("stopCode");
    if (!stop) {
      alert("Please select a bus stop first!");
      return;
    }

    const url = `${serverURL}/routes?stop=${encodeURIComponent(stop)}&code=${encodeURIComponent(code)}`;

   $.get(url, function(data) {
        fillRoutes(data)
        }).fail(function() {
        alert("Failed to fetch routes for the bus stop " + stop);
    });
});

// get arrivals for specific stop and route on button click
$(document).on("click", ".route-button", function() {
    const stop = $("#stop-combobox option:selected").data("stopName");
    const code = $("#stop-combobox option:selected").data("stopCode");
    const route = $(this).val();

    const $container = $("#route-buttons-container");
    $container.find("button").removeClass("btn-primary").addClass("btn-outline-primary"); // remove selection from others
    $(this).removeClass("btn-outline-primary").addClass("btn-primary"); // add to newly selected

    if (!stop) {
      alert("Please select a bus stop first!");
      return;
    }
    if (!route) {
      alert("Please select a route first!");
      return;
    }

    let url;
    if (code) {
        url = `${serverURL}/arrivals?stop=${encodeURIComponent(stop)}&code=${encodeURIComponent(code)}&route=${encodeURIComponent(route)}`;
    } else {
        url = `${serverURL}/arrivals?stop=${encodeURIComponent(stop)}&route=${encodeURIComponent(route)}`;
    }

   $.get(url, function(data) {
        const $container = $("#arrivals-container");
        $container.empty(); // clear container from previous results

        if (data.length == 0) {
            $container.hide(); // hide if no arrivals
        } else {
            $container.show();
            data.forEach(function(arrival) {
                $container.append(
                    $('<span>')
                        .addClass('arrival-span d-block bg-white border rounded p-2 mb-2')
                        .html(`
                            <span class="fw-bold">${arrival.time}</span>, 
                            <span class="text-muted">direction: ${arrival.direction}</span>
                    `)
                );
            });
        }
        }).fail(function() {
        alert("Failed to fetch arrivals for the bus stop and route " + stop + " " + route);
    });
});

$(document).on("click", "#clear-button", function() {
    const $regionCombobox = $("#region-combobox");
    const $stopCombobox = $("#stop-combobox");
    $regionCombobox.val(null);
    $regionCombobox.next(".custom-combobox")
        .find("input")
        .val(null)
    $stopCombobox.empty().val(null);
    $stopCombobox.next(".custom-combobox")
        .find("input")
        .val(null)
    $("#route-buttons-container").empty();
    $("#arrivals-container").empty();
});

// automatic geolocation
function locateUser() {
    if ("geolocation" in navigator) {
        navigator.geolocation.getCurrentPosition(
            function(position) {
                const lat = position.coords.latitude;
                const lon = position.coords.longitude;
                console.log("User coordinates:", lat, lon);

                findNearestStop(lat, lon);
            },
            function(error) {
                console.error("Geolocation error:", error);
            },
            {
                enableHighAccuracy: true,
                timeout: 5000
            }
        );
    } else {
        alert("Geolocation is not supported by your browser.");
    }
}

function selectRegion(region) {
    const $combobox = $("#region-combobox");
    $combobox.val(region);
    $combobox
        .next(".custom-combobox")
        .find("input")
        .val(region);
}

function selectStop(stop) {
    const $combobox = $("#stop-combobox");
    const $option = $combobox.find(`option[value="${stop.id}"]`);
    
    if ($option.length === 0) {
        console.warn("Stop not found:", stop);
        return;
    }
    
    $combobox.val(stop.id);
    const $input = $combobox.next(".custom-combobox").find("input");
    $input.val(`${stop.name}, ${stop.code}`);
    $input.autocomplete("close");
}

function findNearestStop(lat, lon) {
   $.ajax({
    url: `${serverURL}/nearest_stop`,
    method: "POST",
    contentType: "application/json",
    data: JSON.stringify({ lat: "59.44672652334827", lon: "24.894330611091675" }),
    //data: JSON.stringify({ lat: lat.toString(), lon: lon.toString() }),
    dataType: "json",
    success: function(data) {
        selectRegion(data.stop.region);
        fillStops(data.stop.region).done(function() {
                selectStop(data.stop);
        });
        fillRoutes(data.routes);
    }
});
}

