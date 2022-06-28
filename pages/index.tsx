import {useState} from 'react';
import type {NextPage} from 'next'
import SiteFooter from "../components/SiteFooter"
import SiteHead from "../components/SiteHead";
import Form from "../components/Form";

const Home: NextPage = () => {
    return (
        <>
            <SiteHead/>
            <Form/>
            <SiteFooter/>
        </>
    )
}
export default Home
