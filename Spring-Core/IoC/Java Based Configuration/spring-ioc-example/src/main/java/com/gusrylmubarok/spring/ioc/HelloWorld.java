package com.gusrylmubarok.spring.ioc;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class HelloWorld {
    private String message;

    public void setMessage(String message) {
        this.message = message;
    }

    public void getMessage() {
        System.out.println("My Message : " + message);
    }
}