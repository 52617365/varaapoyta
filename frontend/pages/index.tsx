import type {NextPage} from "next";
import Button from "../components/Button";
import React, {useEffect, useState} from "react";
import Card from "../components/Card";
import { ApiResponse } from "../interfaces/api_response_interface";

function Home() {
  const [buttonLoading, setButtonLoading] = React.useState(false);
  const [ravintolatApista, setRavintolat] = React.useState<ApiResponse[]>([]);
  const [is_error, set_error] = React.useState<boolean>(false);
  const [fetched, set_fetched] = React.useState<boolean>(false);
  const [kaupunki, asetaKaupunki] = useState("");
  const handleKaupunki = (event: React.ChangeEvent<HTMLInputElement>) => {
    asetaKaupunki(event.target.value);
  };

  const get_user_city = () => {
    navigator.geolocation.getCurrentPosition(async (position) => {
      try {
        const url = `https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=${position.coords.latitude}&longitude=${position.coords.longitude}&localityLanguage=en`;
        const response = await fetch(url);
        const parsed_response = await response.json();
        asetaKaupunki(parsed_response.city);
      } catch (e) {
      }
    });
  };

  useEffect(() => {
    if (!kaupunki) {
      get_user_city();
    }
  });

  const handleKeypress = async (e: any) => {
    //it triggers by pressing the enter key
    if (e.keyCode === 13) {
      await fetchInformation(kaupunki);
    }
  };
  // TODO: use bigdata geolocation api to get city of user using the thing.
  const fetchInformation = async (city: string) => {
    if (buttonLoading || city == "") {
      return;
    }
    setButtonLoading(true);
    try {
      // const url = `https://www.api.rasmusmaki.com/raflaamo/tables/${city}/1`;
      const url = `http://localhost:8080/tables/${city}`;
      const response = await fetch(url);
      const parsed_response = await response.json();
      setRavintolat(parsed_response);
      setButtonLoading(false);
      set_fetched(true);
      set_error(false);
    } catch (e) {
      setButtonLoading(false);
      set_fetched(true);
      set_error(true);
    }
  };

  const render_results = () => {
    if (ravintolatApista.length === 0 && !fetched) {
      return <></>;
    }

    if (is_error) {
      return <h1>Error fetching endpoint</h1>;
    }
    if (fetched && ravintolatApista.length === 0) {
      return <h1>No restaurants found</h1>;
    }

    // API returns either an error message or an array containing the restaurant information.
    if (Array.isArray(ravintolatApista)) {
      return ravintolatApista.map((restaurant: ApiResponse) => {
        return (
            // Storing the id from the reservation page url as a key so its easy to reuse when in V2 we have reservation too.
            <div key={restaurant.restaurant.reservationPageId}>
              <Card apiResponse={restaurant}/>
            </div>
        );
      });
    } else {
      return <h1>{ravintolatApista}</h1>;
    }
  };
  return (
    <>
      <div className="hero min-h-screen bg-base-200">
        <div className="hero-content text-center">
          <div className="max-w-md pr-5">
            <h1 className="text-xl font-bold">Moikka</h1>
            <p className="py-6">
              Aloita pöytävaraus kirjoittamalla kaupunkisi kenttään.
            </p>
            <div className={"pb-3"}>
              <input
                onKeyDown={handleKeypress}
                type="text"
                placeholder="Kaupunki"
                className="input w-full max-w-xs"
                onChange={handleKaupunki}
                defaultValue={kaupunki}
              />
            </div>
            <div>
              <Button
                text="Hae ravintolat"
                setButton={fetchInformation}
                buttonLoading={buttonLoading}
                textfield_text={kaupunki}
              />
            </div>
            {render_results()}
          </div>
        </div>
      </div>
    </>
  );
}

export default Home;