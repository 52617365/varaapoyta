import { Fragment, useState } from "react";
import { Combobox, Transition } from "@headlessui/react";
import { CheckIcon } from "@heroicons/react/20/solid";

export default function AutoCompletionCity({
  default_value,
}: {
  default_value: string;
}) {
  // TODO: pass this state as argument.
  const [selected, setSelected] = useState(default_value);
  const [query, setQuery] = useState("");

  const filteredPeople =
    query === ""
      ? all_possible_cities
      : all_possible_cities.filter((city) =>
          city.name
            .toLowerCase()
            .replace(/\s+/g, "")
            .includes(query.toLowerCase().replace(/\s+/g, ""))
        );

  return (
    <Combobox value={selected} onChange={setSelected}>
      <div className="relative mt-1 bg-base-200">
        <div className="relative w-full cursor-default overflow-hidden rounded-lg bg-indigo-600 text-left shadow-md focus:outline-none focus-visible:ring-2 focus-visible:ring-black focus-visible:ring-opacity-75 focus-visible:ring-offset-2 focus-visible:ring-offset-bg-indigo-600 sm:text-sm">
          <Combobox.Input
            className="w-full border-none py-2 pl-3 pr-10 text-sm leading-5 text-gray-900 focus:ring-0"
            displayValue={(city: city) => city.name}
            onChange={(event) => setQuery(event.target.value)}
          />
        </div>
        <Transition
          as={Fragment}
          leave="transition ease-in duration-100"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
          afterLeave={() => setQuery("")}
        >
          <Combobox.Options className="absolute mt-1 max-h-60 w-full overflow-auto rounded-md bg-cyan-600 py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm">
            {filteredPeople.length === 0 && query !== "" ? (
              <div className="relative cursor-default select-none py-2 px-4 text-gray-700">
                Nothing found.
              </div>
            ) : (
              filteredPeople.map((person) => (
                <Combobox.Option
                  key={person.id}
                  className={({ active }) =>
                    `relative cursor-default select-none py-2 pl-10 pr-4 ${
                      active ? "text-black" : "text-gray-900"
                    }`
                  }
                  value={person}
                >
                  {({ selected, active }) => (
                    <>
                      <span
                        className={`block truncate ${
                          selected ? "font-medium" : "font-normal"
                        }`}
                      >
                        {person.name}
                      </span>
                      {selected ? (
                        <span
                          className={`absolute inset-y-0 left-0 flex items-center pl-3 ${
                            active ? "text-white" : "text-teal-600"
                          }`}
                        >
                          <CheckIcon className="h-5 w-5" aria-hidden="true" />
                        </span>
                      ) : null}
                    </>
                  )}
                </Combobox.Option>
              ))
            )}
          </Combobox.Options>
        </Transition>
      </div>
    </Combobox>
  );
}

interface city {
  id: number;
  name: string;
}
const all_possible_cities = [
  { id: 1, name: "helsinki" },
  { id: 2, name: "espoo" },
  { id: 3, name: "vantaa" },
  { id: 4, name: "nurmijärvi" },
  { id: 5, name: "kerava" },
  { id: 6, name: "järvenpää" },
  { id: 7, name: "vihti" },
  { id: 8, name: "porvoo" },
  { id: 9, name: "lohja" },
  { id: 10, name: "hyvinkää" },
  { id: 11, name: "karkkila" },
  { id: 12, name: "riihimäki" },
  { id: 13, name: "tallinna" },
  { id: 14, name: "hämeenlinna" },
  { id: 15, name: "lahti" },
  { id: 16, name: "forssa" },
  { id: 17, name: "salo" },
  { id: 18, name: "kotka" },
  { id: 19, name: "kouvola" },
  { id: 20, name: "akaa" },
  { id: 21, name: "loimaa" },
  { id: 22, name: "heinola" },
  { id: 23, name: "hamina" },
  { id: 24, name: "kaarina" },
  { id: 25, name: "turku" },
  { id: 26, name: "kangasala" },
  { id: 27, name: "raisio" },
  { id: 28, name: "tampere" },
  { id: 29, name: "nokia" },
  { id: 30, name: "luumäk" },
  { id: 31, name: "laitila" },
  { id: 32, name: "lappeenranta" },
  { id: 33, name: "mikkeli" },
  { id: 34, name: "rauma" },
  { id: 35, name: "ulvila" },
  { id: 36, name: "pori" },
  { id: 37, name: "jyväskylä" },
  { id: 38, name: "imatra" },
  { id: 39, name: "pieksämäki" },
  { id: 40, name: "savonlinna" },
  { id: 41, name: "varkaus" },
  { id: 42, name: "seinäjoki" },
  { id: 43, name: "kuopio" },
  { id: 44, name: "joensuu" },
  { id: 45, name: "kitee" },
  { id: 46, name: "vaasa" },
  { id: 47, name: "iisalmi" },
  { id: 48, name: "lieksa" },
  { id: 49, name: "kokkola" },
  { id: 50, name: "ylivieska" },
  { id: 51, name: "nurmes" },
  { id: 52, name: "kajaani" },
  { id: 53, name: "sotkamo" },
  { id: 54, name: "muhos" },
  { id: 55, name: "kempele" },
  { id: 56, name: "oulu" },
  { id: 57, name: "rovaniemi" },
  { id: 58, name: "kittilä" },
];
