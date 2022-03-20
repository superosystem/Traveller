package com.gusrylmubarok.springmvc.helloworld.controller;

import com.gusrylmubarok.springmvc.helloworld.model.HelloWorld;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;

import java.time.LocalDateTime;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Controller
public class HelloWorldController {

    public String handler(Model model) {
        HelloWorld helloWorld = new HelloWorld();
        helloWorld.setMessage("HEllo World Example Using Spring MVC");
        helloWorld.setDateTime(LocalDateTime.now().toString());
        model.addAttribute("helloworld", helloWorld);
        return "helloworld";
    }
}
