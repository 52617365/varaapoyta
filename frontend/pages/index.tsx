import type { NextPage } from "next";
import Button from "../components/Button";
import React, { useState } from "react";
import Card from "../components/Card";
import api_response from "../interfaces/api_response_interface";

const Home: NextPage = () => {
  const [buttonLoading, setButtonLoading] = React.useState(false);
  const [ravintolatApista, setRavintolat] = React.useState<api_response[]>([]);
  const [is_error, set_error] = React.useState<boolean>(false);
  const [fetched, set_fetched] = React.useState<boolean>(false);

  const fetchInformation = async (city: string) => {
    if (buttonLoading || city == "") {
      return;
    }
    setButtonLoading(true);
    try {
      const url = `http://localhost:10000/raflaamo/tables/${city}/1`;
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
      const cards = ravintolatApista.map((ravintola: api_response) => {
        return (
          // Storing the id from the reservation page url as a key so its easy to reuse when in V2 we have reservation too.
          <div key={ravintola.links.tableReservationLocalizedId}>
            <Card ravintola={ravintola} />
          </div>
        );
      });
      return cards;
    } else {
      return <h1>{ravintolatApista}</h1>;
    }
  };
  // kaupunki is used in get query to endpoint
  const [kaupunki, asetaKaupunki] = useState("");
  const handleKaupunki = (event: React.ChangeEvent<HTMLInputElement>) => {
    asetaKaupunki(event.target.value);
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
                type="text"
                placeholder="Kaupunki"
                className="input w-full max-w-xs"
                onChange={handleKaupunki}
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
};

export default Home;
