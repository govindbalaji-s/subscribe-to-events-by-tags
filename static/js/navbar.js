/**
 * props.user
 * props.onDashboardNav
 */
function NavBar(props) {
    if(!props.user) {
        return <LoginButton />;
    }
    let loginstatus = (
        <div>
            Hello {props.user[APIUserEmailField]}!
            <br/>
            <LogoutButton />
        </div>
    );
    
    let navlinks = (
        <div>
            <Button onClick={props.onDashboardNav}>
                Home
            </Button>
        </div>
    );
    return (
        <div>
            {loginstatus}
            {navlinks}
        </div>
    );    
}