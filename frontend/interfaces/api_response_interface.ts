export interface ApiResponse {
    restaurant:               Restaurant;
    timeSlots:                string[];
    timeTillRestaurantCloses: TimeTillCloses;
    timeTillKitchenCloses:    TimeTillCloses;
}

export interface Restaurant {
    id:                string;
    reservationPageId: string;
    name:              Name;
    urlPath:           Name;
    address:           Address;
    features:          Features;
    openingTime:       OpeningTime;
    Links:             Links;
}

export interface Links {
    tableReservationLocalized: Name;
    homepageLocalized:         Name;
}

export interface Name {
    fi_FI: string;
}

export interface Address {
    municipality: Name;
    street:       Name;
    zipCode:      string;
}

export interface Features {
    accessible: boolean;
}

export interface OpeningTime {
    restaurantTime: Time;
    kitchenTime:    Time;
}

export interface Time {
    ranges: Range[];
}

export interface Range {
    start:      string;
    end:        string;
    endNextDay: boolean;
}

export interface TimeTillCloses {
    hour:   number;
    minute: number;
}
