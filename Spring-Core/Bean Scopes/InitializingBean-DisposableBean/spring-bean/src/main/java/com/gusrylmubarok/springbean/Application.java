package com.gusrylmubarok.springbean;

import org.springframework.context.annotation.AnnotationConfigApplicationContext;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class Application {
    public static void main(String[] args) {
        AnnotationConfigApplicationContext context = new AnnotationConfigApplicationContext(AppConfig.class);
        context.close();
    }
}
