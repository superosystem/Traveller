package com.gusrylmubarok.springdi;

import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.AnnotationConfigApplicationContext;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class Application {
    public static void main(String[] args) {
        ApplicationContext applicationContext = new AnnotationConfigApplicationContext(SpringDiApplication.class);
        ConstructorBasedInjection  constructorBasedInjection = applicationContext.getBean(ConstructorBasedInjection.class);
        constructorBasedInjection.processMsg("twitter message sending ");
    }
}
