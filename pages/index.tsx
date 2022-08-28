import type {NextPage} from 'next'
import Link from "next/link"
import Image from "next/image"
import SiteFooter from "../components/SiteFooter";
import {useSWRConfig} from "swr";
import Button from "../components/Button";
import React from "react";

const Home: NextPage = () => {
    return (
        <>
            <div className="hero min-h-screen bg-base-200">
                <div className="hero-content text-center">
                    <div className="max-w-md pr-5">
                        <h1 className="text-5xl font-bold">Moikka,</h1>
                        <p className="py-6">Aloita pöytävaraaminen asettamalla mieltymyksesi asetuksissa.</p>
                        <Link href={"/varaa"}>
                            <a className={"absolute top-0 right-0"}>
                                <Button text={"Varaa"}/>
                            </a>
                        </Link>
                        <Link href={"/asetukset"}>
                            <a className={"absolute top-0 right-20"}>
                                <Button text={"Asetukset"}/>
                            </a>
                        </Link>
                        <span>
                            <Link href={"/login"}>
                                <button className="btn btn-square">
                                    <Image alt="gmail logo" src={"/pictures/gmail.svg"} width={50} height={50}/>
                                </button>
                            </Link>
                       </span>
                        <span>
                            <Link href={"/login"}>
                                <button className="btn btn-square">
                                    <Image alt="outlook logo" src={"/pictures/outlook.png"} width={50} height={50}/>
                                </button>
                            </Link>
                        </span>
                        <span>
                                <Link href={"/login"}>
                                    <button className="btn btn-square">
                                        <Image alt="github logo" src={"/pictures/github.svg"} width={50} height={50}/>
                                    </button>
                                </Link>
                        </span>
                    </div>
                </div>
            </div>
            <SiteFooter/>
        </>
    )
}
export default Home
