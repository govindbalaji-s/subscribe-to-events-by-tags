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

/*  Abstract wrapper for a page of a Pagedlist
*/
function PagedListPage(props) {
    let fromIndex = props.fromIndex,
        perPage = props.perPage,
        items = props.items,
        ItemClass = props.itemClass,
        keyFn = props.keyFn,
        otherProps = props.otherProps;
    
    let toIndex = Math.min(fromIndex+perPage, items.length);
    return (
        <div>
            <ul>
                {items.slice(fromIndex, toIndex).map(item =>
                    <ItemClass  key = {keyFn(item)}
                                item = {item}
                                otherProps = {otherProps} />
                    )}
            </ul>
        </div>
    )
}

/*
* An abstract paged list of rows
* Required props:
        props.header => header content before all rows
        props.items  => list of items
        props.itemClass => component class of each item
        props.keyFn     => function that maps each item to a unique key among other props.items
        props.otherProps => other props that each item requires
*/
class PagedList extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            fromIndex: 0,
            perPage: 5
        };
        this.nextPage = this.nextPage.bind(this);
        this.previousPage = this.previousPage.bind(this);
    }

    render() {
        return (
            <div>
                {this.props.header}
                <Button onClick={this.previousPage}>Prev</Button>
                <Button onClick={this.nextPage}>Next</Button>
                <PagedListPage fromIndex = {this.state.fromIndex}
                                perPage = {this.state.perPage}
                                items = {this.props.items}
                                itemClass= {this.props.itemClass}
                                keyFn = {this.props.keyFn}
                                otherProps = {this.props.otherProps}
                />
            </div>
        );
    }

    nextPage() {
        this.setState((state, props) => {
            if(state.fromIndex + state.perPage < props.items.length) {
                return {
                    fromIndex: state.fromIndex+state.perPage
                };
            }
            else{
                return {}
            }
        });
    }

    previousPage() {
        this.setState((state, props) => {
            return {
                fromIndex: Math.max(state.fromIndex-state.perPage, 0)
            };
        });
    }
}

/* Renders a li of a tag, with given tag name, onClick handler and an action prompt
*/
function TagListRow(props) {
    //Required props
    let tag = props.item;
    let onUnfollow = props.otherProps.onUnfollow;
    let actionLabel = props.otherProps.actionLabel;
    let onDetails = props.otherProps.onTagDetails;

    return (
        <li>
            <span>{tag}</span>
            <span>
                <Button onClick={(e) => {onDetails(tag);}}>
                    View
                </Button>
            </span>
            <span>
                <Button onClick={(e) => {onUnfollow(tag);}}>
                    {actionLabel}
                </Button>
            </span>
        </li>
    );
}

/* Renders a li of a event
*/
function EventListRow(props) {
    //Required props
    let event = props.item;
    let onDetails = props.otherProps.onDetails;
    let onAction = props.otherProps.onAction;
    let actionLabel = props.otherProps.actionLabel;

    let name = event[APIEventNameField];
    let time = event[APIEventTimeField];
    let eventID = event[APIEventIDField];

    return (
        <li>
            <span>{name}</span>
            <span>{timestampToString(time)}</span>
            <span>
                <Button onClick={(e) => {onDetails(eventID);}}>
                    View
                </Button>
            </span>
            <span>
                <Button onClick={(e) => {onAction(eventID);}}>
                    {actionLabel}
                </Button>
            </span>
        </li>
    )
}