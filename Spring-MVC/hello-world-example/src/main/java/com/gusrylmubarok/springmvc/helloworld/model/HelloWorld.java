package com.gusrylmubarok.springmvc.helloworld.model;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class HelloWorld {
    private String message;
    private String dateTime;

    public String getMessage() {
        return message;
    }
    public void setMessage(String message) {
        this.message = message;
    }
    public String getDateTime() {
        return dateTime;
    }
    public void setDateTime(String dateTime) {
        this.dateTime = dateTime;
    }
}
