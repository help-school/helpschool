import React from 'react';

//Child components
import Header from "./Header";
import Footer from "./Footer";
//import FeaturedProducts from "./components/Featured";
import School from "./School";
import Introduce from "./Introduce";
import MyDonations from "./MyDonations";

function Home() {

    return (
        <div className="Home">
            {/* <Header/> */}
            <Introduce/>
            <School/>
            {/* <FeaturedProducts/> */}
            <MyDonations/>
        </div>
    );
}

export default Home;
