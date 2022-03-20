package com.gusrylmubarok.springbean;

import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.AnnotationConfigApplicationContext;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class AppSingletonConfigurationTest {
    public static void main(String[] args) {
        ApplicationContext applicationContext = new AnnotationConfigApplicationContext(AppSingletonConfiguration.class);
        UserService  userService = applicationContext.getBean(UserService.class);
        userService.setName("Singleton scope test");
        System.out.println(userService.getName());

        UserService  userService1 = applicationContext.getBean(UserService.class);
        System.out.println(userService1.getName());
    }
}
