import Link from 'next/link'
import Image from 'next/image'

function Footer() {
    return (
        <div className={"fixed left-0 bottom-0 w-full text-center"}>
            <p className={""}>Created by Rase</p>
            <Link className={""} href="https://github.com/52617365">
                <Image src={"/github.png"} alt={"Github Icon"} width={50} height={50}/>
            </Link>
        </div>
    )
}


export default Footer;