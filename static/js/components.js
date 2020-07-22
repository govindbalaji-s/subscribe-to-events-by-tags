function Button(props) {
    return (
        <span className = "button"
             onClick = {e => {props.onClick(e); e.preventDefault();}}
        >
            {props.children}
        </span>
    );
}

//Button that has an href and an onClick
function LinkButton(props) {
    let onClick = (e) => {
        props.onClick && props.onClick(e);
        location.href = props.href;
    };
    return (
        <Button onClick = {onClick}>
            {props.children}
        </Button>
    );
}

function LoginButton(props) {
    return <LinkButton href="/login">Login</LinkButton>;
}

function LogoutButton(props) {
    return (<LinkButton href="/logout" onClick={() => {Cookies.remove('auth-session');}}>
        Logout
    </LinkButton>);
}

/*  Render a row of tags, sliced with given indices
    Renders props.tags[props.fromIndex ... props.fromIndex+props.perPage)
        or till the end of props.tags list

    Each tag has the same props.onClick and props.actionLabel
*/
function TagList(props) {
    //Required props
    let fromIndex = props.fromIndex,
        tags = props.tags,
        onClick = props.onClick,
        actionLabel = props.actionLabel;

    let toIndex = Math.min(props.fromIndex+props.perPage, props.tags.length);
    return (
        <div>
            <ul>
                {tags.slice(fromIndex, toIndex).map( tag => 
                    <TagListRow key = {tag}
                                tag = {tag}
                                onClick = {onClick}
                                actionLabel = {actionLabel}
                    />
                )}
            </ul>
        </div>
    )
}

/* Renders a li of a tag, with given tag name, onClick handler and an action prompt
*/
function TagListRow(props) {
    //Required props
    let name = props.tag;
    let onClick = props.onClick;
    let actionLabel = props.actionLabel;

    return (
        <li>
            <span>{name}</span>
            <span>
                <Button onClick={(e) => {onClick(name);}}>
                    {actionLabel}
                </Button>
            </span>
        </li>
    );
}