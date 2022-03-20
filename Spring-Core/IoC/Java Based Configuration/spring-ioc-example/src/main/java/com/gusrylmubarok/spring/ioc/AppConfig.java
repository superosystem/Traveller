package com.gusrylmubarok.spring.ioc;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Configuration
public class AppConfig {

    @Bean
    public HelloWorld helloWorld() {
        HelloWorld helloWorld = new HelloWorld();
        helloWorld.setMessage("Hello World!");
        return helloWorld;
    }
}
