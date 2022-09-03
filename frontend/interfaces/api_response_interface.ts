interface string_field {
    fi_FI: string
}

interface ranges_contents {
    end: string;
    start: string;
}

interface ranges {
    ranges: Array<ranges_contents>
}

interface opening_time {
    kitchenTime: ranges;
    restaurantTime: ranges;
}

interface address {
    municipality: string_field;
    street: string_field;
    zipCode: string;
}
interface links {
    tableReservationLocalized: string_field;
    tableReservationLocalizedId: string;
    homepageLocalized: string_field;
}

interface api_response {
    id: string;
    name: string_field;
    urlPath: string_field;
    address: address
    openingTime: opening_time;
    links: links;
    available_time_slots: Array<string>;
}

export default api_response;