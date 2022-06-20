package com.gusrylmubarok.sms.controller;

import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;

@Controller
public class HomeController {
    public String homePage() {
        return "index";
    }
}
