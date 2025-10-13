import htmx from "htmx.org"
import $ from "cash-dom"
import axios from "axios"
import { getCookie } from "typescript-cookie"

interface PurchaseResponse {
    checkout_url: string
}

function initLogin(container) {
    $(container).find("#login-form").on("submit", function (e) {
        e.preventDefault()
        axios.post("/login", {
            username: $("#login-username").val(),
            password: $("#login-password").val(),
        }).then(response => {
            const resp = response.data
            window.location.reload()
        }).catch(error => {
            $("#login-error").show().html("Could not login. Please try again.")
        })
    })
}

function makePurchase(orderType: string, itemId: string) {
    const data = {
        order_type: orderType,
        item_id: itemId.toString(),
    }
    axios.post<PurchaseResponse>("/v1/platform/makePurchase", data, {headers: {"X-Session-Id": getCookie("session")}},)
        .then(response => {
            window.location.reload()
        })
}

function initClickBuy(container) {
    $(container).find(".click-pay").on("click", function (e) {
        if (window.appConfig.userId == "") {
            // need to be logged in to buy
            return
        }

        const $e = $(e.target)
        const orderType = $e.data("order-type")
        const itemId = $e.data("item-id")

        if ($e.hasClass("click-pay")) {
            makePurchase(orderType, itemId)
        }
    })
}

function initClose(container) {
    $(container).find(".overlay-widget .close").on("click", function (e) {
        e.preventDefault()
        $(container).html("").hide()
    })
}

function refreshHeader() {
    const onSwap = () => {
        document.body.removeEventListener("htmx:afterSwap", onSwap)
        initLogin(document.body)
    }
    document.body.addEventListener("htmx:afterSwap", onSwap)

    htmx.ajax("get", "/partial/header", {
        target: $("#header-content")[0],
        swap: "innerHTML",
    })
}

refreshHeader()
initClickBuy(document.body)
